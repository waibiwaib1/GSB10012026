"""Extra functions to help preview response content, not used directly by API functions.

These functions will accept any of the following:

* A JSON response
* A list of response objects
* A single response object
"""
from copy import deepcopy
from logging import getLogger
from typing import List, Sequence

from pyinaturalist.constants import ResponseObject, ResponseOrObject

__all__ = [
    'format_controlled_terms',
    'format_identifications',
    'format_observations',
    'format_places',
    'format_projects',
    'format_species_counts',
    'format_taxa',
    'format_users',
    'simplify_observations',
]
logger = getLogger(__name__)


def format_controlled_terms(terms: ResponseOrObject, align: bool = False) -> str:
    """Format controlled term results, with allowed values indented under each term"""
    term_strings = [_format_controlled_term(t, align=align) for t in _ensure_list(terms)]
    return '\n'.join(term_strings)


def _format_controlled_term(term: ResponseObject, align: bool = False) -> str:
    label = f"{term['label']:<30}" if align else term['label']
    term_str = f"[{term['id']:>8}] {label}"
    value_strings = [f"    [{v['id']:>8}] {v['label']}" for v in term.get('values', [])]
    return '\n'.join([term_str] + value_strings)


def format_identifications(identifications: ResponseOrObject, align: bool = False) -> str:
    """Format identification results into a condensed summary: id, taxon, who, and when"""
    id_strings = [_format_identification(i, align=align) for i in _ensure_list(identifications)]
    return '\n'.join(id_strings)


def _format_identification(identification: ResponseObject, align: bool = False) -> str:
    taxon_str = _format_taxon(identification.get('taxon') or {}, align=align)
    user = identification.get('user') or {}
    created_at = identification.get('created_at', '')
    return f"[{identification['id']}] {taxon_str} identified by {user.get('login')} on {created_at}"


def format_observations(observations: ResponseOrObject, align: bool = False) -> str:
    """Format observation results into a condensed summary: id, what, who, when, and where"""
    obs_strings = [_format_observation(o, align=align) for o in _ensure_list(observations)]
    return '\n'.join(obs_strings)


def _format_observation(observation: ResponseObject, align: bool = False) -> str:
    taxon_str = _format_taxon(observation.get('taxon') or {}, align=align)
    location = observation.get('place_guess') or observation.get('location')
    if align:
        return (
            f"[{observation['id']:>8}] {taxon_str}"
            f"\n    observed by {observation['user']['login']} "
            f"on {observation['observed_on']} at {location}"
        )
    else:
        return (
            f"[{observation['id']}] {taxon_str} "
            f"observed by {observation['user']['login']} "
            f"on {observation['observed_on']} at {location}"
        )


def format_places(places: ResponseOrObject, align: bool = False) -> str:
    """Format place results into a single string containing place ID and display name"""
    place_strings = [_format_place(p, align=align) for p in _ensure_list(places)]
    return '\n'.join(place_strings)


def _format_place(place: ResponseObject, align: bool = False) -> str:
    name = place.get('display_name') or place.get('name')
    if align:
        return f"[{place['id']:>8}] {name}"
    else:
        return f"[{place['id']}] {name}"


def format_projects(projects: ResponseOrObject, align: bool = False) -> str:
    """Format project results into a single string containing project ID and title"""
    project_strings = [_format_project(p, align=align) for p in _ensure_list(projects)]
    return '\n'.join(project_strings)


def _format_project(project: ResponseObject, align: bool = False) -> str:
    if align:
        return f"[{project['id']:>8}] {project['title']}"
    else:
        return f"[{project['id']}] {project['title']}"


def format_species_counts(species_counts: ResponseOrObject, align: bool = False) -> str:
    """Format observation species counts"""
    count_strings = [_format_species_count(t, align=align) for t in _ensure_list(species_counts)]
    return '\n'.join(count_strings)


def _format_species_count(species_count: ResponseObject, align: bool = False) -> str:
    taxon = _format_taxon(species_count['taxon'], align=align)
    return f'{taxon}: {species_count["count"]}'


def format_taxa(taxa: ResponseOrObject, align: bool = False) -> str:
    """Format taxon results into a single string containing taxon ID, rank, and name
    (including common name, if available).
    """
    taxon_strings = [_format_taxon(t, align=align) for t in _ensure_list(taxa)]
    return '\n'.join(taxon_strings)


def _format_taxon(taxon: ResponseObject, align: bool = False) -> str:
    if not taxon:
        return 'unknown taxon'
    common_name = taxon.get('preferred_common_name')
    name = taxon.get('name') or 'unknown'
    name += f' ({common_name})' if common_name else ''
    rank = taxon.get('rank', '').title()

    # Visually align taxon IDs (< 7 chars) and ranks (< 11 chars)
    if align:
        return f"[{taxon['id']:>8}] {rank:>12} {name:>55}"
    elif rank:
        return f"[{taxon['id']}] {rank}: {name}"
    else:
        return f"[{taxon['id']}] {name}"


def format_users(users: ResponseOrObject, align: bool = False) -> str:
    """Format user results into a single string containing user ID, login, and real name
    (if available).
    """
    user_strings = [_format_user(u, align=align) for u in _ensure_list(users)]
    return '\n'.join(user_strings)


def _format_user(user: ResponseObject, align: bool = False) -> str:
    name = f"{user['login']}" + (f" ({user['name']})" if user.get('name') else '')
    if align:
        return f"[{user['id']:>8}] {name}"
    else:
        return f"[{user['id']}] {name}"


def simplify_observations(
    observations: ResponseOrObject, align: bool = False
) -> List[ResponseObject]:
    """Flatten out some nested data structures within observation recorda:

    * annotations
    * comments
    * identifications
    * non-owner IDs
    """
    return [_simplify_observation(o) for o in _ensure_list(observations)]


def _simplify_observation(obs):
    # Reduce annotations to IDs and values
    obs = deepcopy(obs)
    obs['annotations'] = [
        (a['controlled_attribute_id'], a['controlled_value_id']) for a in obs['annotations']
    ]

    # Reduce identifications to just a list of identification IDs and taxon IDs
    obs['identifications'] = [(i['id'], i['taxon_id']) for i in obs['identifications']]
    obs['non_owner_ids'] = [(i['id'], i['taxon_id']) for i in obs['non_owner_ids']]

    # Reduce comments to usernames and comment text
    obs['comments'] = [(c['user']['login'], c['body']) for c in obs['comments']]
    del obs['observation_photos']

    return obs


def _ensure_list(obj: ResponseOrObject) -> List:
    if isinstance(obj, dict) and 'results' in obj:
        return obj['results']
    elif isinstance(obj, Sequence):
        return list(obj)
    else:
        return [obj]
